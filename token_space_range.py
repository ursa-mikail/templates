import math

def max_value(base, digits):
    """Calculate max decimal value that can be represented with given base and digit count."""
    return base**digits - 1

def storage_size(base, digits, char_bit_size):
    """Estimate total storage size in bits for a number represented in base N."""
    return digits * char_bit_size

# Example Cases
base_n_cases = [(2, 4, 8), (16, 4, 8), (64, 4, 8), (64, 6, 8)]  # (Base, Digits, Bits per symbol)

for base, digits, bits in base_n_cases:
    max_val = max_value(base, digits)
    storage = storage_size(base, digits, bits)
    print(f"Base {base}, Digits {digits}: Max Value = {max_val}, Storage = {storage} bits, Value range sizing = {max_val + 1}")


import string

BASE64_ALPHABET = string.ascii_uppercase + string.ascii_lowercase + string.digits + "+/"

def base10_to_base64(n):
    if n == 0:
        return "A"  # Base64 encoding of zero
    result = ""
    while n > 0:
        result = BASE64_ALPHABET[n % 64] + result
        n //= 64
    return result

print(base10_to_base64(1000))  # Example
print(f"len of symbols for base64: {len(BASE64_ALPHABET)}")


"""
Base 2, Digits 4: Max Value = 15, Storage = 32 bits, Value range sizing = 16
Base 16, Digits 4: Max Value = 65535, Storage = 32 bits, Value range sizing = 65536
Base 64, Digits 4: Max Value = 16777215, Storage = 32 bits, Value range sizing = 16777216
Base 64, Digits 6: Max Value = 68719476735, Storage = 48 bits, Value range sizing = 68719476736
Po
len of symbols for base64: 64
"""