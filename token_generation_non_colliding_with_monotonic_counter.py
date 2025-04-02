import pyffx
import string
import matplotlib.pyplot as plt
from collections import Counter

# Define alphabet sets
BASE64_ALPHABET = string.ascii_uppercase + string.ascii_lowercase + string.digits + "+/"

def base10_to_baseN(n, alphabet, length=5):
    """Convert a base-10 number to a base-N string with fixed length."""
    base = len(alphabet)
    result = ""
    while n > 0:
        result = alphabet[n % base] + result
        n //= base
    result = result.rjust(length, alphabet[0])  # Pad with the first character to ensure length
    return result[-length:]  # Ensure fixed length

def baseN_to_base10(s, alphabet):
    """Convert a base-N string back to base-10 integer."""
    base = len(alphabet)
    n = 0
    for char in s:
        n = n * base + alphabet.index(char)
    return n

def apply_fpe_and_plot_distribution(limit, alphabet, key, length=5):
    """Increment from 0 to limit, apply FPE, convert back to integer, and plot distribution."""
    e = pyffx.String(key.encode(), alphabet=alphabet, length=length)  # Fixed-length encryption
    encrypted_values = [e.encrypt(base10_to_baseN(i, alphabet, length)) for i in range(limit)]
    decrypted_values = [baseN_to_base10(val, alphabet) for val in encrypted_values]

    #print(f"encrypted_values: {encrypted_values}")
    #print(f"decrypted_values: {decrypted_values}")
    
    # Plot distribution
    plt.figure(figsize=(10, 5))
    plt.hist(decrypted_values, bins=50, color='blue', alpha=0.7)
    plt.xlabel("Decrypted Integer Value")
    plt.ylabel("Frequency (within the range)")
    plt.title("Distribution of Decrypted FPE Values")
    plt.show()
    
    print(f"Total unique encrypted values: {len(set(encrypted_values))} out of {limit}")

# Parameters
M = 10000000  # Large number limit
KEY = "secret-key"
ALPHABET = BASE64_ALPHABET  # Base64 Alphabet

apply_fpe_and_plot_distribution(M, ALPHABET, KEY)

"""

"""