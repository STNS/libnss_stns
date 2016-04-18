require 'spec_helper'

describe file("/etc/stns/libnss_stns.conf") do
  it { should be_mode(644) }
  it { should be_owned_by('root') }
  it { should be_grouped_into('root') }
end

files = if os[:family] == 'redhat'
  if i386?
    %w(
      /usr/lib/libnss_stns.so
      /lib/security/libpam_stns.so
      /usr/lib/libnss_stns.so.2
    )
  else
    %w(
      /usr/lib64/libnss_stns.so
      /lib64/security/libpam_stns.so
      /usr/lib64/libnss_stns.so.2
    )
  end
elsif ['debian', 'ubuntu'].include?(os[:family])
  if i386?
    %w(
      /usr/lib/i386-linux-gnu/libnss_stns.so
      /lib/i386-linux-gnu/security/libpam_stns.so
      /lib/i386-linux-gnu/libnss_stns.so.2
    )
  else
    %w(
      /usr/lib/x86_64-linux-gnu/libnss_stns.so
      /lib/x86_64-linux-gnu/security/libpam_stns.so
      /lib/x86_64-linux-gnu/libnss_stns.so.2
    )
  end
end

files.each_with_index do |f,i|
  describe file(f) do
    it { should be_owned_by('root') }
    it { should be_grouped_into('root') }
    it { should be_symlink } if i == 2
  end

  describe command("file #{f}") do
    bit = i386? ? "32" : "64"
    its(:stdout) { should match /#{bit}-bit/ }
  end
end
