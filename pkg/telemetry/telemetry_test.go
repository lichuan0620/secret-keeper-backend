package telemetry

import (
	"context"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

const tcTimeout = 3

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Telemetry")
}

var _ = Describe("Server", func() {
	const listenAddress = "127.0.0.1:9090"
	const healthzAddress = "http://" + listenAddress + healthzEndpoint
	const metricsAddress = "http://" + listenAddress + metricsEndpoint

	var server *Server
	var ctx context.Context
	var cancel context.CancelFunc
	var startErrCh chan error

	BeforeEach(func() {
		server = NewServer(&ServerOptions{
			ListenAddress: listenAddress,
		})
		ctx, cancel = context.WithTimeout(context.Background(), tcTimeout*time.Second)
		startErrCh = make(chan error)
	})

	AfterEach(func() {
		cancel()
		Expect(<-startErrCh).To(Succeed())
	})

	Describe("healthz", func() {

		It("can return OK for a passed check", func() {
			go func() {
				startErrCh <- server.Start(ctx)
			}()
			Eventually(func(g Gomega) {
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthzAddress, nil)
				g.Expect(err).ToNot(HaveOccurred())
				resp, err := http.DefaultClient.Do(req)
				g.Expect(err).ToNot(HaveOccurred())
				defer func() {
					_ = resp.Body.Close()
				}()
				g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
			}).Should(Succeed())
		})

		It("can respond with an error for a failed check", func() {
			server.SetHealthCheck("fail", func(_ *http.Request) error {
				return errors.New("pseudo error")
			})
			go func() {
				startErrCh <- server.Start(ctx)
			}()
			Eventually(func(g Gomega) {
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthzAddress, nil)
				g.Expect(err).ToNot(HaveOccurred())
				resp, err := http.DefaultClient.Do(req)
				g.Expect(err).ToNot(HaveOccurred())
				defer func() {
					_ = resp.Body.Close()
				}()
				g.Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			}).Should(Succeed())
		})
	})
})
