#!/usr/bin/env python3

# Copyright 2020 NLP Odyssey Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import torch

FLOAT_DTYPES = [
    torch.float16,  # or torch.half
    torch.float32,  # or torch.float
    torch.float64,  # or torch.double
    torch.bfloat16
]

INT_DTYPES = [
    torch.int8,
    torch.int16,  # or torch.short
    torch.int32,  # or torch.int
    torch.int64,  # or torch.long
]


def main():
    for proto in range(1, 6):
        for use_zip in [False, True]:
            for dtype in FLOAT_DTYPES:
                save([1.2, -3.4, 5.6, -7.8], dtype, proto, use_zip)

            for dtype in INT_DTYPES:
                save([1, -2, 3, -4], dtype, proto, use_zip)

            save([1, 10, 100, 255], torch.uint8, proto, use_zip)
            save([True, False, True, False], torch.bool, proto, use_zip)


def save(data, dtype, proto, use_zip):
    str_dtype = str(dtype)[6:]
    str_zip = '_zip' if use_zip else ''
    torch.save(
        torch.tensor(data, dtype=dtype),
        f'tensor_{str_dtype}_proto{proto}{str_zip}.pt',
        pickle_protocol=proto,
        _use_new_zipfile_serialization=use_zip)


if __name__ == '__main__':
    main()
