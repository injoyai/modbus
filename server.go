package modbus

import (
	"bufio"
	"context"
	"fmt"
	"github.com/goburrow/serial"
	"io"
	"log"
	"net"
)

type RTUOption = serial.Config

type Server struct {
	Coils            Coils                //线圈 0x01(0x读),0x05(写1个),0x0f(写多个)
	DiscreteInputs   Coils                //离散输入(只读的线圈) 0x02(1x读)
	InputRegisters   Register             //输入寄存器 0x04(3x读)
	HoldingRegisters Register             //保持寄存器 0x03(4x读) 0x06(写1个) 0x10(写多个)
	Handler          [43]Handler          //处理函数,下标对应功能码
	listenTCP        []net.Listener       //监听tcp
	ctxTCP           context.Context      //TCP上下文
	cancelTCP        context.CancelFunc   //TCP上下文关闭
	listenRTU        []io.ReadWriteCloser //监听rtu
	ctxRTU           context.Context      //RTU上下文
	cancelRTU        context.CancelFunc   //RTU上下文关闭
	ctx              context.Context      //上下文
	debug            bool                 //打印日志
	printHandler     func(Frame, Frame)   //打印日志函数
}

// SetCoils 设置线圈接口
func (this *Server) SetCoils(register uint16, wrc ReadWriteCoils) {
	this.Coils[register] = wrc
}

// SetDiscreteInputs 设置离散输入接口
func (this *Server) SetDiscreteInputs(register uint16, wrc ReadWriteCoils) {
	this.DiscreteInputs[register] = wrc
}

// SetInputRegisters 设置数据寄存器接口
func (this *Server) SetInputRegisters(register uint16, wrc ReadWriteRegister) {
	this.InputRegisters[register] = wrc
}

// SetHoldingRegisters 设置保持寄存器接口
func (this *Server) SetHoldingRegisters(register uint16, wrc ReadWriteRegister) {
	this.HoldingRegisters[register] = wrc
}

// SetHandler 设置功能码对应函数
func (this *Server) SetHandler(code int, handler Handler) {
	if code > 0 && code < len(this.Handler) {
		this.Handler[code] = handler
	}
}

// SetPrintHandler 设置打印通讯日志函数
func (this *Server) SetPrintHandler(fn func(origin, result Frame)) *Server {
	this.printHandler = fn
	return this
}

// Debug 调试模式
func (this *Server) Debug(b ...bool) {
	this.debug = !(len(b) > 0 && !b[0])
}

// ListenTCPNum 监听TCP端口的服务数量
func (this *Server) ListenTCPNum() int {
	return len(this.listenTCP)
}

// ListenRTUNum 监听RTU服务的数量
func (this *Server) ListenRTUNum() int {
	return len(this.listenRTU)
}

// Close 关闭所有
func (this *Server) Close() error {
	this.CloseRTU()
	this.CloseTCP()
	return nil
}

// CloseRTU 关闭所有RTU连接
func (this *Server) CloseRTU() {
	if this.cancelRTU != nil {
		this.cancelRTU()
		this.cancelRTU = nil
	}
	this.listenRTU = []io.ReadWriteCloser(nil)
}

// CloseTCP 关闭所有TCP连接
func (this *Server) CloseTCP() {
	if this.cancelTCP != nil {
		this.cancelTCP()
		this.cancelTCP = nil
	}
	this.listenTCP = []net.Listener(nil)
}

// ListenTCP 监听TCP端口
func (this *Server) ListenTCP(port int) error {
	if this.cancelTCP == nil {
		this.ctxTCP, this.cancelTCP = context.WithCancel(this.ctx)
	}
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	this.listenTCP = append(this.listenTCP, listen)
	go func(ctx context.Context, listen net.Listener) {
		defer listen.Close()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				conn, err := listen.Accept()
				if err != nil {
					this.printErr(err)
					return
				}
				go func(ctx context.Context, conn net.Conn) {
					defer conn.Close()
					for {
						select {
						case <-ctx.Done():
							return
						default:
							//按TCP格式读取数据
							frame, err := ReadWithTCP(conn)
							if err != nil {
								this.printErr(err)
								return
							}
							this.printErr(this.handle(frame, conn))
						}
					}
				}(ctx, conn)
			}
		}
	}(this.ctxTCP, listen)
	return nil
}

// ListenRTU 监听RTU
func (this *Server) ListenRTU(cfg *serial.Config) error {
	if this.cancelRTU == nil {
		this.ctxRTU, this.cancelRTU = context.WithCancel(this.ctx)
	}
	client, err := serial.Open(cfg)
	if err != nil {
		return err
	}
	this.listenRTU = append(this.listenRTU, client)
	go func(ctx context.Context, conn serial.Port) {
		defer conn.Close()
		buf := bufio.NewReader(conn)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				//按RTU读取数据
				frame, err := ReadWithRTU(buf)
				if err != nil {
					this.printErr(err)
					continue
				}
				this.printErr(this.handle(frame, conn))
			}
		}
	}(this.ctxRTU, client)
	return nil
}

func (this *Server) printErr(err error) {
	if this.debug && err != nil {
		log.Println("[错误]", err)
	}
}

func NewServer() *Server {
	return NewServerWithContext(context.Background())
}

func NewServerWithContext(ctx context.Context) *Server {
	s := &Server{}
	s.ctx = ctx
	s.SetPrintHandler(s.defaultPrintHandler)
	s.SetHandler(1, s.handler1)
	s.SetHandler(2, s.handler2)
	s.SetHandler(3, s.handler3)
	s.SetHandler(4, s.handler4)
	s.SetHandler(5, s.handler5)
	s.SetHandler(6, s.handler6)
	s.SetHandler(15, s.handler15)
	s.SetHandler(16, s.handler16)
	return s
}
