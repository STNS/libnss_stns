task :default => "build"

desc "build binary"
task "build" do
  sh "docker build --no-cache --rm -t stns:stns ."
  sh "docker run -v \"$(pwd)\"/binary:/go/src/github.com/STNS/libnss_stns/binary -t stns:stns"
end

desc "build rpm"
task "rpm" => [:build] do
  docker_run "rpm"
end

desc "build deb"
task "deb" => [:build] do
  docker_run "deb"
end

def docker_run(file, dir="binary")
  sh "docker build --no-cache --rm -f docker/#{file} -t stns:stns ."
  sh "docker run -it -v \"$(pwd)\"/binary:/go/src/github.com/STNS/libnss_stns/#{dir} -t stns:stns"
end
