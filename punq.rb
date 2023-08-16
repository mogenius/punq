class Punq < Formula
  desc "View your kubernetes workloads relativly neat!"
  homepage "https://punq-k8s.io"
  
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/mogenius/homebrew-punq/releases/download/v1.0.6/punq-1.0.6-darwin-arm64.tar.gz"
      sha256 "07cc2c38eeea1504eb44bf0aa216dbb6d23768a21b2770dc7582fdaf295a69cd"
    elsif Hardware::CPU.intel?
      url "https://github.com/mogenius/homebrew-punq/releases/download/v1.0.6/punq-1.0.6-darwin-amd64.tar.gz"
      sha256 ""
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/mogenius/homebrew-punq/releases/download/v1.0.6/punq-1.0.6-linux-amd64.tar.gz"
        sha256 ""
      else
        url "https://github.com/mogenius/homebrew-punq/releases/download/v1.0.6/punq-1.0.6-linux-386.tar.gz"
        sha256 ""
      end
    elif Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/mogenius/homebrew-punq/releases/download/v1.0.6/punq-1.0.6-linux-arm64.tar.gz"
        sha256 ""
      else
        url "https://github.com/mogenius/homebrew-punq/releases/download/v1.0.6/punq-1.0.6-linux-arm.tar.gz"
        sha256 ""
      end
    end
  end
  
  version "1.0.6"
  license "MIT"

  def install
  if OS.mac?
    if Hardware::CPU.arm?
      # Installation steps for macOS ARM64
      bin.install "punq-1.0.6-darwin-arm64" => "punq"
    elsif Hardware::CPU.intel?
      # Installation steps for macOS AMD64
      bin.install "punq-1.0.6-darwin-amd64" => "punq"
    end
  elsif OS.linux?
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        # Installation steps for Linux AMD64
        bin.install "punq-1.0.6-linux-amd64" => "punq"
      else
        # Installation steps for Linux 386
        bin.install "punq-1.0.6-linux-386" => "punq"
      end
    elsif Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        # Installation steps for Linux ARM64
        bin.install "punq-1.0.6-linux-arm64" => "punq"
      else
        # Installation steps for Linux ARM
        bin.install "punq-1.0.6-linux-arm" => "punq"
      end
    end
  end
end
end
