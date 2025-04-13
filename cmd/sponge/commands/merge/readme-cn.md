## 合并代码命令

根据 Protobuf 文件生成的新代码将自动合并到本地 Go 文件中，并确保不会影响已有的业务逻辑。如果合并过程中发生异常，系统会在 /tmp/sponge_merge_backup_code 目录中备份合并前的代码，方便您随时恢复。

### 使用说明

1. 基于 protobuf 创建的 web 服务

   ```bash
   sponge merge http-pb --dir=serverDir
   ```

2. 基于 protobuf 或 sql 创建的 gRPC 服务

   ```bash
   sponge merge rpc-pb --dir=serverDir
   ```
3. 基于 protobuf 创建的 grpc 网关服务

   ```bash
   sponge merge rpc-gw-pb --dir=serverDir
   ```
4. 基于 protobuf 创建的 gRPC+HTTP 混合服务

   ```bash
   sponge merge grpc-http-pb --dir=serverDir
   ```