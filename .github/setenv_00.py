#!/usr/bin/env python3

# fmt: off

import os
import sys
from typing import NoReturn


_basedir = os.path.abspath(os.path.dirname(__file__))
_basedir = os.path.abspath(os.path.join(_basedir, '..'))

def _print(msg: str):
    print(msg, file=sys.stderr)
def show_errmsg(errmsg: str) -> NoReturn:
    _print(f'[e] {errmsg}')
    sys.exit(1)


if __name__ == "__main__":
    _github_env = os.getenv('GITHUB_ENV', '')
    if not _github_env:
        show_errmsg('This script should be run on Github Action')

    _pkg_plat = sys.argv[1]
    _pkg_arch = sys.argv[2]
    _pkg_libc = sys.argv[3]

    def _setenv(f, k, v):
        _print(f'{k}: {v}'); f.write(f'{k}={v}\n')

    with open(_github_env, 'a') as f:
        _inst_dir = os.path.abspath(os.path.join(_basedir, 'out'))
        _setenv(f, 'INST_DIR', _inst_dir)
        _deps_ver = f'{_inst_dir}.ver'
        _setenv(f, 'DEPS_VER', _deps_ver)

        _target_arch_libc = f'{_pkg_arch}'
        if _pkg_plat == 'linux' and _pkg_libc:
            _target_arch_libc = f'{_pkg_arch}-{_pkg_libc}'
        _setenv(f, 'TARGET_ARCH_LIBC', _target_arch_libc)
