import os
import subprocess


class GenProto:
    modName = "first"
    servicePath = "./service"

    def __init__(self, workDir):
        self.workDir = workDir

    def check(self):
        with os.scandir(os.path.join(self.workDir, "idl")) as it:
            for entry in it:
                if entry.name.endswith(".proto") and entry.is_file():
                    self.done(entry.name)

    def makedirs(self, path) -> bool:
        try:
            print(path)
            os.stat(path)
            return True
        except:
            os.makedirs(path)
            return False

    def existF(self, path) -> bool:
        try:
            os.stat(path)
            return True
        except:
            return False


    def done(self, filename):
        service = os.path.join(self.servicePath, filename.rsplit('.')[0])

        cmd = "kitex -module %s  -type protobuf -service user -I ./idl  ./idl/%s" % (self.modName, filename.__str__())

        cmd += " && rm build.sh kitex.yaml ./script -rf"

        self.makedirs(service)
        exist = self.existF(os.path.join(service, "main.go"))
        if exist:
            cmd += " && rm main.go handler.go"
        else:
            cmd += " && mv main.go handler.go %s" % (service)

        try:
            print("[exec]: ", cmd)
            subprocess.run(cmd, shell=True, check=True)
        except:
            print("some thing error")
            subprocess.run("exit 1", shell=True)

        print("[done] ", filename)


if __name__ == '__main__':
    a = GenProto(".")
    a.check()
