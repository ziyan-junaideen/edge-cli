class EdgeCli < Formula
  desc "Command line client for Edge Payment Technologies"
  homepage "https://github.com/ziyan-junaideen/edge-cli"
  url "https://github.com/ziyan-junaideen/edge-cli.git", tag: "v0.1.0"
  license "MIT"
  head "https://github.com/ziyan-junaideen/edge-cli.git", branch: "main"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X github.com/ziyan-junaideen/edge-cli/internal/cli.version=#{version}"
    system "go", "build", *std_go_args(ldflags: ldflags, output: bin/"edge"), "./cmd/edge"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/edge --version")
    assert_match "Edge Payment Technologies", shell_output("#{bin}/edge --help")
  end
end
