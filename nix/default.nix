{ lib, buildGoModule, version ? "unstable" }:

buildGoModule {
  pname = "warss";
  inherit version;

  src = lib.fileset.toSource {
    root = ../.;
    fileset = lib.fileset.intersection (lib.fileset.fromSource (lib.sources.cleanSource ../.)) (
      lib.fileset.unions [
      ../go.mod
      ../go.sum
      ../main.go
      ../internal
      ]
    );
  };

  vendorHash = null;

  ldflags = [
    "-s"
    "-w"
    "-X main.version=${version}"
  ];

  meta = {
    description = "A TUI RSS feed";
    homepage = "https://github.com/pixel-87/warss";
    license = lib.licenses.gpl3Plus;
    maintainers = with lib.maintainers; [ pixel-87 ];
    mainProgram = "warss";
  };
}
