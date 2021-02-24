# docker-sync

同步 Docker 镜像到指定镜像仓库

`config.yml`配置如下：

> `[]`为可选

```yaml
sync:
  - src: REPOSITORY[:TAG][@DIGEST]
    dst: REPOSITORY[:TAG]
```

同步后输出`synced.yml`结构:

```yaml
sync:
  - src: REPOSITORY:TAG@DIGEST
    dst: REPOSITORY:TAG@DIGEST
```

## 过程

### Pull Image

```shell script
$ docker pull hello-world:latest@sha256:1a523af650137b8accdaed439c17d684df61ee4d74feac151b5b337bd29e7eec

sha256:1a523af650137b8accdaed439c17d684df61ee4d74feac151b5b337bd29e7eec: Pulling from library/hello-world
Digest: sha256:1a523af650137b8accdaed439c17d684df61ee4d74feac151b5b337bd29e7eec
Status: Image is up to date for hello-world@sha256:1a523af650137b8accdaed439c17d684df61ee4d74feac151b5b337bd29e7eec
docker.io/library/hello-world:latest@sha256:1a523af650137b8accdaed439c17d684df61ee4d74feac151b5b337bd29e7eec
```

### 通过 REPOSITORY:TAG 查看 ImageId

```shell script
$ docker images --format {{.ID}} hello-world:latest

bf756fb1ae65
```

### 为 ImageId 打 Tag，并 Push

```shell script
$ docker tag bf756fb1ae65 registry.cn-hangzhou.aliyuncs.com/hb-chen/hello-world:v0.0.1
```

```shell script
$ docker push registry.cn-hangzhou.aliyuncs.com/hb-chen/hello-world:v0.0.1

The push refers to repository [registry.cn-hangzhou.aliyuncs.com/hb-chen/hello-world]
9c27e219663c: Preparing
9c27e219663c: Layer already exists
v0.0.1: digest: sha256:90659bf80b44ce6be8234e6ff90a1ac34acbeb826903b02cfa0da11c82cbc042 size: 525
```

Push 后获取 Digest,输出同步后 Image 列表
