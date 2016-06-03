require 'erb'
require 'rake'
require 'rspec/core/rake_task'

task :default => "test"

desc "run server spec"
task :spec    => 'spec:all'

desc "run unit test"
task "test" do
  sh 'docker run --rm -it -v "$(pwd)":/go/src/github.com/STNS/libnss_stns -t pyama/stns:ubuntu-x86 /bin/bash'
end

desc "make package all architecture"
task "make_pkg" => %W(
  clean_pkg
  make_pkg_x86
  make_pkg_i386
)

desc "test package all architecture"
task "test_pkg" => %W(
  make_pkg
  test_pkg_x86
  test_pkg_i386
)

%w(x86 i386).each do |arch|
  desc "make package #{arch}"
  task "make_pkg_#{arch}" => %W(
    clean_bin
    ubuntu_build_#{arch}
    centos_pkg_#{arch}
    ubuntu_pkg_#{arch}
  )

  desc "test package #{arch}"
  task "test_pkg_#{arch}" => %W(
    centos_ci_#{arch}
    ubuntu_ci_#{arch}
  )
end

task "clean_all" do
  sh "rm -rf binary/*"
end

desc "delete binarys"
task "clean_bin" do
  sh "find binary/* | grep -v -e 'rpm$' -e 'deb$' | xargs rm -rf"
end

desc "delete packages"
task "clean_pkg" do
  sh "find binary/* | grep -e 'rpm$' -e 'deb$' | xargs rm -rf"
end

[
  {
    os: "centos",
    arch: %w(x86 i386),
    pkg_arch: %w(x86_64 i386)
  },
  {
    os: "ubuntu",
    arch: %w(x86 i386),
    pkg_arch: %w(amd64 i386)
  }
].each do |h|
  h[:arch].each_with_index do |arch,index|
    unless h[:os] == "centos"
      desc "#{h[:os]}_build_#{arch}"
      task "#{h[:os]}_build_#{arch}" do
        docker_run(h[:os], arch, "build")
      end
    end

    %w(pkg ci).each do |t|
      task "#{h[:os]}_#{t}_#{arch}" do
        docker_run(h[:os], arch, t, h[:pkg_arch][index])
      end
    end
  end
end

def docker_run(os, arch, task, pkg_arch=nil)
  content = ERB.new(open("docker/#{os}-#{task}.erb").read).result(binding)
  open("docker/tmp/#{os}-#{arch}-#{task}","w") {
    |f| f.write(content)
  }

  sh "docker build --rm -f docker/tmp/#{os}-#{arch}-#{task} -t stns:lib_stns ."
  sh "docker run --rm -it -v \"$(pwd)\"/binary:/go/src/github.com/STNS/libnss_stns/binary -t stns:lib_stns"
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
