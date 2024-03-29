package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"reflect"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
	"github.com/growerlab/growerlab/src/go-git-grpc/server/command"
)

type File struct {
	Ref         string
	AuthorName  string
	AuthorEmail string
	Message     string
	FilePath    string
	FileContent []byte
}

func (f *File) Verify() error {
	switch true {
	case govalidator.IsNull(f.Ref),
		govalidator.IsNull(f.AuthorName),
		govalidator.IsNull(f.AuthorEmail),
		govalidator.IsNull(f.Message),
		govalidator.IsNull(f.FilePath),
		len(f.FileContent) == 0:
		return errors.New("invalid fields, all is required.")
	}
	return nil
}

type Door struct {
	ctx    context.Context //
	client pb.DoorClient   //
}

func NewDoor(ctx context.Context, pbClient pb.DoorClient) *Door {
	door := &Door{
		ctx:    ctx,
		client: pbClient,
	}
	return door
}

func (d *Door) AddOrUpdateFile(ctx *command.Context, file *File) (string, error) {
	switch true {
	case ctx == nil,
		file == nil,
		file.Verify() != nil:
		panic(errors.New("invalid params"))
	}

	resp, err := d.client.AddOrUpdateFile(d.ctx, &pb.AddFileRequest{
		Path:        ctx.RepoPath,
		Bin:         ctx.Bin,
		Ref:         file.Ref,
		AuthorName:  file.AuthorName,
		AuthorEmail: file.AuthorEmail,
		Message:     file.Message,
		FilePath:    file.FilePath,
		FileContent: file.FileContent,
	})
	if err != nil {
		return "", errors.Trace(err)
	}
	return resp.CommitHash, nil
}

func (d *Door) RunCommand(params *command.Context) error {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("ServeUploadPack panic: %+v", e)
		}
	}()

	cmdResult, err := d.client.RunCommand(d.ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = d.sendContextPack(cmdResult, params); err != nil {
		return err
	}

	return d.copy(cmdResult, params.In, params.Out)
}

func (d *Door) copy(pipe clientStream, in io.Reader, out io.Writer) (err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		var endReason = "normal"
		defer wg.Done()
		defer log.Printf("scan is done, reason '%s'.\n", endReason)

		if in == nil || reflect.ValueOf(in).IsNil() {
			return
		}

		var (
			inReader = bufio.NewReader(in)
			buf      = make([]byte, 32*1024)
		)

		for {
			select {
			case <-d.ctx.Done():
				endReason = "ctx done"
				goto END
			default:
				n := 0
				n, err = inReader.Read(buf)
				if err != nil {
					if err != io.EOF {
						endReason = fmt.Sprintf("read err: %+v", err)
					}
					goto END
				}
				if n <= 0 {
					continue
				}
				err = pipe.Send(&pb.Request{Raw: buf[:n]})
				if err != nil {
					endReason = fmt.Sprintf("send err: %+v", err)
					goto END
				}
			}
		}
	END:
		_ = pipe.CloseSend()
	}()

	wg.Add(1)
	go func() {
		var endReason = "normal"
		defer wg.Done()
		defer log.Printf("recv is done, reason '%s'.\n", endReason)

		for {
			select {
			case <-d.ctx.Done():
				endReason = "ctx done"
				goto END
			default:
				var resp *pb.Response
				resp, err = pipe.Recv()
				if err != nil {
					if err != io.EOF {
						endReason = fmt.Sprintf("recv err: %+v", err)
					}
					goto END
				}
				_, err = out.Write(resp.Raw)
				if err != nil {
					endReason = fmt.Sprintf("write err: %+v", err)
					goto END
				}
			}
		}
	END:
	}()

	wg.Wait()

	if err == io.EOF {
		err = nil // ignore
	}
	return
}

type clientStream interface {
	CloseSend() error
	Send(*pb.Request) error
	Recv() (*pb.Response, error)
}

func (d *Door) sendContextPack(pack clientStream, params *command.Context) error {
	firstReq := &pb.Request{
		Path:      params.RepoPath,
		Env:       params.Env,
		Bin:       params.Bin,
		Args:      params.Args,
		Deadline:  uint64(params.Deadline),
		Raw:       nil,
		HasInput:  params.In != nil,
		HasOutput: params.Out != nil,
	}
	err := pack.Send(firstReq)
	return errors.WithStack(err)
}
