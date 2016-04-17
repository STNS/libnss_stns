require 'erb'
task :default => "test"

task "make" => [:pkg_x86, :pkg_i386]

task "test" do
  docker_run("ubuntu-x86-test")
end

task "clean_bin" do
  sh "find binary/* | grep -v -e 'rpm$' -e 'deb$' | xargs rm -rf"
end

[
  %w(x86 x86_64 amd64),
  %w(i386 i386 i386)
].each do |r|
  arch = r[0]
  arch_rpm = r[1]
  arch_deb = r[2]

  task "build_#{arch}" => [:clean_bin]  do
    content = ERB.new(open("docker/ubuntu-build.erb").read).result(binding)
    open("docker/tmp/ubuntu-#{arch}-build","w") {
      |f| f.write(content)
    }
    docker_run "tmp/ubuntu-#{arch}-build"
  end

  task "pkg_#{arch}" => ["build_#{arch}".to_sym] do
    [
      ["centos", arch_rpm, "rpm"],
      ["ubuntu", arch_deb, "deb"]
    ].each do |o|
      content = ERB.new(open("docker/#{o[0]}-pkg.erb").read).result(binding)
      open("docker/tmp/#{o[0]}-#{arch}-pkg","w") {
        |f| f.write(content)
      }

      sh "find binary/* | grep -e '#{o[1]}.#{o[2]}$' | xargs rm -rf"

      docker_run("tmp/#{o[0]}-#{arch}-pkg", o[1])
      # check package
      sh "test -e binary/libnss*#{o[1]}.#{o[2]}"
      sh "test -e binary/libpam*#{o[1]}.#{o[2]}"
    end
  end
end

def docker_run(file, arch="x86_64", dir="binary")
  sh "docker build --no-cache --rm -f docker/#{file} -t stns:stns ."
  sh "docker run -e ARCH=#{arch} --rm -it -v \"$(pwd)\"/binary:/go/src/github.com/STNS/libnss_stns/#{dir} -t stns:stns"
end
