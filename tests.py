import os
import subprocess
import zlib

import pandas as pd
import tqdm


def zopfli_compress(sample_set: str) -> int:
    with subprocess.Popen(
        [
            "zopfli",
            "-c",
            "--deflate",
            f"examples/{sample_set}",
        ],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    ) as process:
        process.wait()
        if process.returncode != 0:
            raise RuntimeError("Failed to compress")

        return len(process.stdout.read())


def zlib_compress(sample_set: str) -> int:
    with open(f"examples/{sample_set}", "rb") as file:
        return len(zlib.compress(file.read(), 9))


def deflate_compress(sample_set: str) -> int:
    with subprocess.Popen(
        [
            "./deflate",
            "-c",
            "-in",
            f"examples/{sample_set}",
        ],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    ) as process:
        process.wait()
        if process.returncode != 0:
            raise RuntimeError(
                process.stderr.read().decode().strip().removeprefix("Error: ")
            )

        return len(process.stdout.read())


if __name__ == "__main__":
    sets = [
        "README.md",
        "alice29.txt",
        "alphabet.txt",
        "asyoulik.txt",
        "cp.html",
        "fields.c",
        "grammar.lsp",
        "helloworld.txt",
        "random.txt",
        "sum",
        "xargs.1",
    ]

    data = []
    for sample_set in tqdm.tqdm(sets):
        path = f"examples/{sample_set}"
        size_before = os.path.getsize(path)

        row = [sample_set]
        for compress in [zopfli_compress, zlib_compress, deflate_compress]:
            size_after = compress(sample_set)
            row.append(f"{size_after / size_before * 100:.1f}%")

        data.append(row)

    df = pd.DataFrame(data, columns=["Set", "Zopfli", "ZLIB", "Deflate"])
    df.to_csv("tests.csv", index=False)
