package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"reflect"
	"sync"

	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/go-git-grpc/pb"
	"github.com/growerlab/growerlab/src/go-git-grpc/server/git"
)

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

func (d *Door) RunGit(params *git.Context) error {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("ServeUploadPack panic: %+v", e)
		}
	}()

	runGit, err := d.client.RunGit(d.ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = d.sendContextPack(runGit, params); err != nil {
		return err
	}

	return d.copy(runGit, params.In, params.Out)
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

func (d *Door) sendContextPack(pack clientStream, params *git.Context) error {
	firstReq := &pb.Request{
		Path:      params.RepoPath,
		Env:       params.Env,
		GitBin:    params.GitBin,
		Args:      params.Args,
		Deadline:  uint64(params.Deadline),
		Raw:       nil,
		HasInput:  params.In != nil,
		HasOutput: params.Out != nil,
	}
	err := pack.Send(firstReq)
	return errors.WithStack(err)
}
