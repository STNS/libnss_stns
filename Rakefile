task :default => "build"

desc "test"
task "test" do
  docker_run("test")
end

desc "build binary 64bit"
task "build_64" do
  sh "docker build --no-cache --rm -t stns:stns ."
  sh "docker run -v \"$(pwd)\"/binary:/go/src/github.com/STNS/libnss_stns/binary -t stns:stns"
end

desc "make package 64bit"
task "pkg_64" => [:build_64] do
  docker_run("rpm", "x86_64")
  docker_run("deb", "amd64")
end

desc "make binary 32bit"
task "build_32" do
  docker_run "build_32"
end

desc "make package 32bit"
task "pkg_32" => [:build_32] do
  docker_run("rpm", "i386")
  docker_run("deb_32", "i386")
end

def docker_run(file, arch="x86_64", dir="binary")
  sh "docker build --no-cache --rm -f docker/#{file} -t stns:stns ."
  sh "docker run -e TARGET=#{arch} -it -v \"$(pwd)\"/binary:/go/src/github.com/STNS/libnss_stns/#{dir} -t stns:stns"
end
