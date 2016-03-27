task :default => "test"

task "make" => [:pkg_x86, :pkg_i386]

task "test" do
  docker_run("ubuntu-x86-test")
end

task "clean_bin" do
  sh "ls -d binary/* | grep -v -e 'rpm$' -e 'deb$' | xargs rm -rf"
end

[
  %w(x86 x86_64 amd64),
  %w(i386 i386 i386)
].each do |r|
  task "build_#{r[0]}" => [:clean_bin]  do
    docker_run "ubuntu-#{r[0]}-build"
  end

  task "pkg_#{r[0]}" => ["build_#{r[0]}".to_sym] do
    sh "ls -d binary/* | grep -e '#{r[1]}.rpm$' -e '#{r[2]}.deb$'| xargs rm -rf"
    docker_run("centos-#{r[0]}-rpm", r[1])
    docker_run("ubuntu-#{r[0]}-deb", r[2])

    # check package
    sh "test -e binary/*#{r[1]}.rpm"
    sh "test -e binary/*#{r[2]}.deb"
  end
end

def docker_run(file, arch="x86_64", dir="binary")
  sh "docker build --no-cache --rm -f docker/#{file} -t stns:stns ."
  sh "docker run -e ARCH=#{arch} -it -v \"$(pwd)\"/binary:/go/src/github.com/STNS/lib-stns/#{dir} -t stns:stns"
end
