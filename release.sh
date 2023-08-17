#!/bin/bash

BINARY_NAME=punq
VERSION=$(grep "Ver" version/consts.go | awk -F "\"" {'print $2'})
SHA256_DARWIN_ARM64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-darwin-arm64.tar.gz | awk '{print $1}')
SHA256_DARWIN_AMD64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-darwin-amd64.tar.gz | awk '{print $1}')
SHA256_LINUX_ARM64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-linux-arm64.tar.gz | awk '{print $1}')
SHA256_LINUX_ARM=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-linux-arm.tar.gz | awk '{print $1}')
SHA256_LINUX_AMD64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-linux-amd64.tar.gz | awk '{print $1}')
SHA256_LINUX_386=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-linux-386.tar.gz | awk '{print $1}')

# Generate formula from template with replacements
cat <<EOF > punq.rb
class Punq < Formula
  desc "View your kubernetes workloads relativly neat!"
  homepage "https://punq.dev"
  
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/mogenius/punq/releases/download/v${VERSION}/punq-${VERSION}-darwin-arm64.tar.gz"
      sha256 "$SHA256_DARWIN_ARM64"
    elsif Hardware::CPU.intel?
      url "https://github.com/mogenius/punq/releases/download/v${VERSION}/punq-${VERSION}-darwin-amd64.tar.gz"
      sha256 "$SHA256_DARWIN_AMD64"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/mogenius/punq/releases/download/v${VERSION}/punq-${VERSION}-linux-amd64.tar.gz"
        sha256 "$SHA256_LINUX_AMD64"
      else
        url "https://github.com/mogenius/punq/releases/download/v${VERSION}/punq-${VERSION}-linux-386.tar.gz"
        sha256 "$SHA256_LINUX_386"
      end
    elif Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/mogenius/punq/releases/download/v${VERSION}/punq-${VERSION}-linux-arm64.tar.gz"
        sha256 "$SHA256_LINUX_ARM64"
      else
        url "https://github.com/mogenius/punq/releases/download/v${VERSION}/punq-${VERSION}-linux-arm.tar.gz"
        sha256 "$SHA256_LINUX_ARM"
      end
    end
  end
  
  version "${VERSION}"
  license "MIT"

  def install
  if OS.mac?
    if Hardware::CPU.arm?
      # Installation steps for macOS ARM64
      bin.install "punq-$VERSION-darwin-arm64" => "punq"
    elsif Hardware::CPU.intel?
      # Installation steps for macOS AMD64
      bin.install "punq-$VERSION-darwin-amd64" => "punq"
    end
  elsif OS.linux?
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        # Installation steps for Linux AMD64
        bin.install "punq-$VERSION-linux-amd64" => "punq"
      else
        # Installation steps for Linux 386
        bin.install "punq-$VERSION-linux-386" => "punq"
      end
    elsif Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        # Installation steps for Linux ARM64
        bin.install "punq-$VERSION-linux-arm64" => "punq"
      else
        # Installation steps for Linux ARM
        bin.install "punq-$VERSION-linux-arm" => "punq"
      end
    end
  end
end
end
EOF
