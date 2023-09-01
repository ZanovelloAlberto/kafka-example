
{ config, pkgs, requireFile, ... }:

self: super:

{
  go = super.go.overrideAttrs ( old:  {
    name = "go";
    version = "1.21.0";
    }); 
}
