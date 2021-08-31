import docker
import sys

name = sys.argv[1].split("=")[1]
image = sys.argv[2].split("=")[1]
ports = {}
environments = {}
volumes = {'/root/logs': {'bind': '/mnt/cuack/', 'mode': 'rw'}}

for x in sys.argv[3:]:
    if x.startswith("ports"):
        x = x[6:]
        x = x.split(",")
        for p in x[1:]:
            ports[int(p)] = int(p)
    else:
        x = x.split("=")
        if "-" in x[0]:
            x[0] = x[0].replace('-', '')
        environments[(x[0]).upper()] = x[1]

print(name)
print(image)  
print(environments)
print(ports)

client = docker.from_env()
container = client.containers.run(image=image, detach=True, name=name, ports=ports, environment=environments, volumes=volumes)