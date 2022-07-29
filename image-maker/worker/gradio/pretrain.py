# encoding: utf-8

from collections import namedtuple
import json
import sys


_Pretrain = namedtuple("Pretrain", ["owner", "repo", "file"])


class _InvalidPretrain(Exception):
    pass


def _parse_pretrain(path):
    if path == "":
        return None

    s = path.strip().strip("/")
    v = s.split("/")
    if len(v) < 3:
        raise(_InvalidPretrain("invalid pretrain path"))

    return _Pretrain(v[0], v[1], s)


def _load_config(path):
    with open(path, 'r') as f:
        data = json.load(f)

        return _parse_pretrain(data.get("pretrain"))


def load(path):
    v = _load_config(path)

    if v is None:
        return

    print("%s\n%s\n%s" % v)


if __name__ == "__main__":
    if len(sys.argv) != 2:
        sys.exit(1)

    try:
        load(sys.argv[1])
    except Exception as e:
        print(e)
        sys.exit(1)
