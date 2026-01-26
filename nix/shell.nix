{
  mkShell,
  callPackage,

  go,
  gopls,
  gofumpt,
  goreleaser,
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
  ];
}
