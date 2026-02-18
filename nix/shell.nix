{
  mkShell,
  callPackage,

  go_1_26,
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
    go_1_26
    gopls
    gofumpt
    goreleaser
    golangci-lint
  ];
}
