## English | [简体中文](readme-cn.md)

## Merge Code Commands

The newly generated code based on the Protobuf file will be automatically merged into the local Go file, ensuring that existing business logic remains unaffected. If an error occurs during the merge process, a backup of the pre-merge code will be stored in the `/tmp/sponge_merge_backup_code` directory, allowing you to restore the previous code state if needed.

### Usage Instructions

1. Create web server based on protobuf

   ```bash
   sponge merge http-pb --dir=serverDir
   ```

2.  Create gRPC server based on protobuf or sql

   ```bash
   sponge merge rpc-pb --dir=serverDir
   ```

3. Create gRPC gateway server based on protobuf

   ```bash
   sponge merge rpc-gw-pb --dir=serverDir
   ```

4. Create gRPC + HTTP server based on protobuf

   ```bash
   sponge merge grpc-http-pb --dir=serverDir
   ```
