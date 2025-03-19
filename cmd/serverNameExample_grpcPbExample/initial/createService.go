package initial

import (
	"strconv"

	"github.com/go-dev-frame/sponge/pkg/app"

	"github.com/go-dev-frame/sponge/internal/config"
	"github.com/go-dev-frame/sponge/internal/server"
)

// CreateServices create grpc service
func CreateServices() []app.IServer {
	var cfg = config.Get()
	var servers []app.IServer
	var grpcAddr = ":" + strconv.Itoa(cfg.Grpc.Port)

	// case 1, create a grpc service without registry
	grpcServer := server.NewGRPCServer(grpcAddr)

	// case 2, create a grpc service and register it with consul or etcd or nacos
	//grpcRegistry, grpcInstance := registerService("grpc", cfg.App.Host, cfg.Grpc.Port)
	//grpcServer := server.NewGRPCServer(grpcAddr,
	//	server.WithGrpcRegistry(grpcRegistry, grpcInstance),
	//)

	servers = append(servers, grpcServer)

	return servers
}

// register service with consul or etcd or nacos, select one of them to use
//func registerService(scheme string, host string, port int) (registry.Registry, *registry.ServiceInstance) {
//	var (
//		instanceEndpoint = fmt.Sprintf("%s://%s:%d", scheme, host, port)
//		cfg              = config.Get()
//
//		iRegistry registry.Registry
//		instance  *registry.ServiceInstance
//		err       error
//
//		id       = cfg.App.Name + "_" + scheme + "_" + host + "_" + strconv.Itoa(port)
//		logField logger.Field
//	)
//
//	switch cfg.App.RegistryDiscoveryType {
//	case "consul":
//		iRegistry, instance, err = consul.NewRegistry(
//			cfg.Consul.Addr,
//			id,
//			cfg.App.Name,
//			[]string{instanceEndpoint},
//		)
//		if err != nil {
//			panic(err)
//		}
//		logField = logger.Any("consulAddress", cfg.Consul.Addr)
//
//	case "etcd":
//		iRegistry, instance, err = etcd.NewRegistry(
//			cfg.Etcd.Addrs,
//			id,
//			cfg.App.Name,
//			[]string{instanceEndpoint},
//		)
//		if err != nil {
//			panic(err)
//		}
//		logField = logger.Any("etcdAddress", cfg.Etcd.Addrs)
//
//	case "nacos":
//		iRegistry, instance, err = nacos.NewRegistry(
//			cfg.NacosRd.IPAddr,
//			cfg.NacosRd.Port,
//			cfg.NacosRd.NamespaceID,
//			id,
//			cfg.App.Name,
//			[]string{instanceEndpoint},
//		)
//		if err != nil {
//			panic(err)
//		}
//		logField = logger.String("nacosAddress", fmt.Sprintf("%v:%d", cfg.NacosRd.IPAddr, cfg.NacosRd.Port))
//	}
//
//	if instance != nil {
//		msg := fmt.Sprintf("register service address to %s", cfg.App.RegistryDiscoveryType)
//		logger.Info(msg, logger.String("name", cfg.App.Name), logger.String("endpoint", instanceEndpoint), logger.String("id", id), logField)
//		return iRegistry, instance
//	}
//
//	return nil, nil
//}
