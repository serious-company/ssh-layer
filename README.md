# SSH Layer builder

To build a docker image with your `id_rsa` tha is locate on `~/.ssh` you can run the command:
```shell
ssh-layer build
```
You need docker version 19 with api version 1.40, after you run it an image with the name `ssh-layer` will be
on your images `docker images | grep ssh-layer`
