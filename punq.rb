class Punq < Formula
  desc "View your kubernetes workloads relativly neat!"
  homepage "https://punq-k8s.io"
  
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/mogenius/punq/releases/download/v1.0.7/punq-1.0.7-darwin-arm64.tar.gz"
      sha256 "f836ec92e7cfbeb03aa37085aa14fa76eb3a53180bf6e425ce39219716c21529"
    elsif Hardware::CPU.intel?
      url "https://github.com/mogenius/punq/releases/download/v1.0.7/punq-1.0.7-darwin-amd64.tar.gz"
      sha256 "66f052f289d24099ecfa3b4af97552069d4222aeadbacd55a5b8450cb6d14bcc"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/mogenius/punq/releases/download/v1.0.7/punq-1.0.7-linux-amd64.tar.gz"
        sha256 "32a243ca0f7f1201790b83d4ef3d801e22b1812a5f0aaa288e8db45df5c748a0"
      else
        url "https://github.com/mogenius/punq/releases/download/v1.0.7/punq-1.0.7-linux-386.tar.gz"
        sha256 "dc7c7549ffc093917b6decc0a3fce3de96c2a86ef1dd3d18e3039afe61176c6c"
      end
    elif Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/mogenius/punq/releases/download/v1.0.7/punq-1.0.7-linux-arm64.tar.gz"
        sha256 "2790cc4c0f806a8a7dd7e7b58be1683b06af935d944dd10e05a207fad476b804"
      else
        url "https://github.com/mogenius/punq/releases/download/v1.0.7/punq-1.0.7-linux-arm.tar.gz"
        sha256 "c237222f83c8d32fa66cc41af32943043fbcb31a91974e2fa27f1d8563815cdd"
      end
    end
  end
  
  version "1.0.7"
  license "MIT"

  def install
  if OS.mac?
    if Hardware::CPU.arm?
      # Installation steps for macOS ARM64
      bin.install "punq-1.0.7-darwin-arm64" => "punq"
    elsif Hardware::CPU.intel?
      # Installation steps for macOS AMD64
      bin.install "punq-1.0.7-darwin-amd64" => "punq"
    end
  elsif OS.linux?
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        # Installation steps for Linux AMD64
        bin.install "punq-1.0.7-linux-amd64" => "punq"
      else
        # Installation steps for Linux 386
        bin.install "punq-1.0.7-linux-386" => "punq"
      end
    elsif Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        # Installation steps for Linux ARM64
        bin.install "punq-1.0.7-linux-arm64" => "punq"
      else
        # Installation steps for Linux ARM
        bin.install "punq-1.0.7-linux-arm" => "punq"
      end
    end
  end
end
end
