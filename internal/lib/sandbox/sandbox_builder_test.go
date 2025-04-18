package sandbox_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	types "k8s.io/cri-api/pkg/apis/runtime/v1"

	"github.com/cri-o/cri-o/internal/hostport"
	libsandbox "github.com/cri-o/cri-o/internal/lib/sandbox"
)

var _ = Describe("Sandbox:Builder", func() {
	BeforeEach(func() {
		builder = libsandbox.NewBuilder()
		builder.SetCreatedAt(time.Now())
		builder.SetLogDir("test")
		builder.SetShmPath("test")
		builder.SetNamespace("")
		builder.SetKubeName("")
		builder.SetMountLabel("test")
		builder.SetProcessLabel("test")
		builder.SetCgroupParent("")
		builder.SetPrivileged(false)
		builder.SetRuntimeHandler("")
		builder.SetResolvPath("")
		builder.SetHostname("")
		builder.SetPortMappings([]*hostport.PortMapping{})
		builder.SetHostNetwork(false)
		builder.SetUsernsMode("")
		builder.SetPodLinuxOverhead(nil)
		builder.SetPodLinuxResources(nil)
		err := builder.SetCRISandbox(builder.ID(), make(map[string]string), make(map[string]string), &types.PodSandboxMetadata{})
		Expect(err).ToNot(HaveOccurred())
	})

	t.Describe("SetConfig", func() {
		It("should succeed", func() {
			// Given
			config := &types.PodSandboxConfig{
				Metadata: &types.PodSandboxMetadata{
					Name:      "name",
					Uid:       "uid",
					Namespace: "namespace",
				},
			}
			Expect(builder.SetConfig(config)).To(Succeed())

			// When
			err := builder.GenerateNameAndID()

			// Then
			Expect(err).ToNot(HaveOccurred())
			Expect(builder.ID()).To(HaveLen(64))
			Expect(builder.Name()).To(ContainSubstring("name"))
			Expect(builder.Name()).To(ContainSubstring("uid"))
			Expect(builder.Name()).To(ContainSubstring("namespace"))

			_, err = builder.GetSandbox()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with empty config", func() {
			// Given
			// When
			err := builder.GenerateNameAndID()

			// Then
			Expect(err).To(HaveOccurred())
			Expect(builder.ID()).To(Equal(""))
			Expect(builder.Name()).To(Equal(""))
		})

		It("should fail with empty name in metadata", func() {
			// Given
			config := &types.PodSandboxConfig{
				Metadata: &types.PodSandboxMetadata{
					Uid:       "uid",
					Namespace: "namespace",
				},
			}
			Expect(builder.SetConfig(config)).NotTo(Succeed())

			// When
			err := builder.GenerateNameAndID()

			// Then
			Expect(err).To(HaveOccurred())
			Expect(builder.ID()).To(Equal(""))
			Expect(builder.Name()).To(Equal(""))
		})

		It("should fail with empty namespace in metadata", func() {
			// Given
			config := &types.PodSandboxConfig{
				Metadata: &types.PodSandboxMetadata{
					Name: "name",
					Uid:  "uid",
				},
			}
			Expect(builder.SetConfig(config)).To(Succeed())

			// When
			err := builder.GenerateNameAndID()

			// Then
			Expect(err).To(HaveOccurred())
			Expect(builder.ID()).To(Equal(""))
			Expect(builder.Name()).To(Equal(""))
		})

		It("should fail with empty uid in metadata", func() {
			// Given
			config := &types.PodSandboxConfig{
				Metadata: &types.PodSandboxMetadata{
					Name:      "name",
					Namespace: "namespace",
				},
			}
			Expect(builder.SetConfig(config)).To(Succeed())

			// When
			err := builder.GenerateNameAndID()

			// Then
			Expect(err).To(HaveOccurred())
			Expect(builder.ID()).To(Equal(""))
			Expect(builder.Name()).To(Equal(""))
		})
	})
	t.Describe("SetConfig", func() {
		BeforeEach(func() {
			builder.SetID("id")
			builder.SetName("")
		})

		It("should succeed", func() {
			// Given
			config := &types.PodSandboxConfig{
				Metadata: &types.PodSandboxMetadata{Name: "name"},
			}
			// When
			err := builder.SetConfig(config)

			// Then
			Expect(err).ToNot(HaveOccurred())
			Expect(builder.Config()).To(Equal(config))
			_, err = builder.GetSandbox()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with nil config", func() {
			// Given
			// When
			err := builder.SetConfig(nil)

			// Then
			Expect(err).To(HaveOccurred())
			Expect(builder.Config()).To(BeNil())

			_, err = builder.GetSandbox()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with empty config", func() {
			// Given
			config := &types.PodSandboxConfig{}

			// When
			err := builder.SetConfig(config)

			// Then
			Expect(err).To(HaveOccurred())
			Expect(builder.Config()).To(BeNil())

			_, err = builder.GetSandbox()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with an empty name", func() {
			// Given
			config := &types.PodSandboxConfig{
				Metadata: &types.PodSandboxMetadata{},
			}

			// When
			err := builder.SetConfig(config)

			// Then
			Expect(err).To(HaveOccurred())
			Expect(builder.Config()).To(BeNil())

			_, err = builder.GetSandbox()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with config already set", func() {
			// Given
			config := &types.PodSandboxConfig{
				Metadata: &types.PodSandboxMetadata{Name: "name"},
			}
			err := builder.SetConfig(config)
			Expect(err).ToNot(HaveOccurred())

			// When
			err = builder.SetConfig(config)

			// Then
			Expect(err).To(HaveOccurred())
			Expect(builder.Config()).NotTo(BeNil())

			_, err = builder.GetSandbox()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
