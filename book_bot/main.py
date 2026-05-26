from stats import get_num_words, get_char_count, get_char_count_sorted
import sys

def get_book_text(filepath):
    file_contents = ""
    with open(filepath) as f:
        file_contents = f.read()
    return file_contents


def main():
    # book_path = "books/frankenstein.txt"
    if len(sys.argv) != 2:
        print("Usage: python3 main.py <path_to_book>")
        sys.exit(1)
    # print(sys.argv)
    book_path = sys.argv[1]
    text = get_book_text(book_path)
    num_words = get_num_words(text)
    # print(f"Found {num_words} total words")
    char_counts = get_char_count(text)
    # print(char_counts)
    sorted_char_counts = get_char_count_sorted(text)
    # print(sorted_char_counts)

    print("============ BOOKBOT ============")
    print(f"Analyzing book found at {book_path}...")
    print("----------- Word Count ----------")
    print(f"Found {num_words} total words")
    print("--------- Character Count -------")
    for items in sorted_char_counts:
        print(f"{items['char']}: {items['count']}")
    print("============= END ===============")

main()

