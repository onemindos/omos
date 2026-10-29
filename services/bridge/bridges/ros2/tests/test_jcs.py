import json
import os

import pytest

from ros2bridge import jcs

VEC = json.load(open(os.path.join(os.path.dirname(__file__), "fixtures", "policy_vectors.json")))


@pytest.mark.parametrize("v,want", [
    # RFC 8785 Appendix B / section 3.2.2.3 number examples
    (0.0, "0"), (-0.0, "0"), (1.0, "1"), (4.5, "4.5"), (0.002, "0.002"), (1e-7, "1e-7"),
    (1e21, "1e+21"), (1e20, "100000000000000000000"), (123456789012345680000.0, "123456789012345680000"),
    (5e-324, "5e-324"), (1.7976931348623157e308, "1.7976931348623157e+308"), (333333333.3333333, "333333333.3333333"),
    (-1.5e-7, "-1.5e-7"), (0.000001, "0.000001"), (9007199254740992, "9007199254740992"), (30.2672, "30.2672"),
])
def test_numbers(v, want):
    assert jcs.canonicalize(v).decode() == want


def test_rfc8785_key_sort_and_strings():
    # RFC 8785 section 3.2.3 example: sorted by UTF-16 code units
    obj = {"€": "Euro Sign", "\r": "Carriage Return", "דּ": "Hebrew Letter Dalet With Dagesh",
           "1": "One", "\U0001f600": "Emoji: Grinning Face", "\u0080": "Control", "ö": "Latin Small Letter O With Diaeresis"}
    keys = list(json.loads(jcs.canonicalize(obj)).keys())
    assert keys == ["\r", "1", "\u0080", "ö", "€", "\U0001f600", "דּ"]
    assert jcs.canonicalize("a\"\\\x01\n é") == '"a\\"\\\\\\u0001\\n é"'.encode()
    assert jcs.canonicalize({"b": [1, True, None], "a": {}}) == b'{"a":{},"b":[1,true,null]}'


def test_nan_refused():
    with pytest.raises(ValueError):
        jcs.canonicalize(float("nan"))


@pytest.mark.parametrize("vec", VEC["vectors"], ids=lambda v: v["name"])
def test_matches_policy_service_canonical_form(vec):
    env = dict(vec["envelope"])
    env.pop("sig")
    assert jcs.canonicalize(env).decode() == vec["jcs_without_sig"]
