with import <nixpkgs> {}; rec {
  cplateEnv = stdenv.mkDerivation {
    name = "cmake";
    buildInputs = [ stdenv
                    cmake
                    pkgconfig
                    gcc
                    gdb
                    clang
                    cquery
                    clang-tools
                    llvm
                    llvmPackages.libclang
                    lldb
                    ctags
                    cppcheck
                    kcov
                    valgrind
                    xorg.libX11
                    libyamlcpp
                  ];
    LD_LIBRARY_PATH="/run/opengl-driver/lib;${xorg.libX11}/lib/;${libGL}/lib/;${libGLU}/lib;${freeglut}/lib";
  };
}
