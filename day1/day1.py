from pathlib import Path

def get_line_value(line: str) -> int:
    numbers = [char for char in line if char.isdigit()]
    return int(numbers[0] + numbers[-1])


def get_sum(filename: str) -> int:
    result = 0
    with Path(filename).open() as f:
        for line in f.readlines():
            result += get_line_value(line)

    return result


# test get_line_value
assert get_line_value("1abc2") == 12
assert get_line_value("pqr3stu8vwx") == 38
assert get_line_value("a1b2c3d4e5f") == 15
assert get_line_value("treb7uchet") == 77

# test get_sum
test_data = "1abc2\npqr3stu8vwx\na1b2c3d4e5f\ntreb7uchet"
test_file = Path("day_1_test_data.txt")
test_file.write_text(test_data)
assert get_sum("day_1_test.txt") == 142
test_file.unlink()
