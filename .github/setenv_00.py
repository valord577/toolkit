#!/usr/bin/env python3

# fmt: off

import sys
from pathlib import Path

sys.dont_write_bytecode = True
sys.path.append(
    (Path(__file__).parent / '..').absolute().as_posix()
)
from scripts import utils as x
# ----------------------------

import os


_github_env = os.getenv('GITHUB_ENV')
if not _github_env:
    raise RuntimeError('This script should be run on Github Action')

_pkg_plat = sys.argv[1]
_pkg_arch = sys.argv[2]
_pkg_libc = sys.argv[3]



with open(_github_env, 'a') as f:
    _target_arch_libc = f'{_pkg_arch}'
    if _pkg_plat in ['linux', 'android'] and _pkg_libc:
        _target_arch_libc = f'{_pkg_arch}-{_pkg_libc}'
    x._util_append_ci_env(f, 'TARGET_ARCH_LIBC', _target_arch_libc)
