{
  mkShell,
  callPackage,

  go,
  gopls,
  gofumpt,
  goreleaser,
  golangci-lint,
}:

let
  defaultPackage = callPackage ./default.nix { };
in
mkShell {
  inputsFrom = [ defaultPackage ];

  packages = [
    go
    gopls
    gofumpt
    goreleaser
    golangci-lint
  ];
}
