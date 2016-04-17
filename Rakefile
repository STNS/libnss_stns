require 'erb'
require 'rake'
require 'rspec/core/rake_task'

task :default => "test"

task :spec    => 'spec:all'
task "make" => [:pkg_x86, :pkg_i386]
task "pkg_test" => [:ci_x86, :ci_i386]

task "test" do
  docker_run "ubuntu-x86-test"
end

task "clean_bin" do
  sh "find binary/* | grep -v -e 'rpm$' -e 'deb$' | xargs rm -rf"
end

[
  %w(x86 x86_64 amd64),
  %w(i386 i386 i386)
].each do |arch, arch_rpm, arch_deb|
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
    ].each do |os, pkg_arch, pkg|
      content = ERB.new(open("docker/#{os}-pkg.erb").read).result(binding)
      open("docker/tmp/#{os}-#{arch}-pkg","w") {
        |f| f.write(content)
      }

      sh "find binary/* | grep -e '#{pkg_arch}.#{pkg}$' | xargs rm -rf"
      docker_run("tmp/#{os}-#{arch}-pkg", pkg_arch)
    end
  end

  task "ci_#{arch}" => ["pkg_#{arch}"] do
    [
      ["centos", arch_rpm],
      ["ubuntu", arch_deb]
    ].each do |os, pkg_arch|
      content = ERB.new(open("docker/#{os}-ci.erb").read).result(binding)
      open("docker/tmp/#{os}-#{arch}-ci","w") {
        |f| f.write(content)
      }

      docker_run("tmp/#{os}-#{arch}-ci", pkg_arch)
    end
  end
end

def docker_run(file, arch="x86_64", dir="binary")
  sh "docker build --rm -f docker/#{file} -t stns:stns ."
  sh "docker run -e ARCH=#{arch} --rm -it -v \"$(pwd)\"/binary:/go/src/github.com/STNS/libnss_stns/#{dir} -t stns:stns"
end

namespace :spec do
  targets = []
  Dir.glob('./spec/*').each do |dir|
    next unless File.directory?(dir)
    target = File.basename(dir)
    target = "_#{target}" if target == "default"
    targets << target
  end

  task :all     => targets
  task :default => :all

  targets.each do |target|
    original_target = target == "_default" ? target[1..-1] : target
    desc "Run serverspec tests to #{original_target}"
    RSpec::Core::RakeTask.new(target.to_sym) do |t|
      ENV['TARGET_HOST'] = original_target
      t.pattern = "spec/#{original_target}/*_spec.rb"
    end
  end
end
