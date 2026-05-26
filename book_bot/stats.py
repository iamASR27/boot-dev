from typing import TypedDict

class CharacterCount(TypedDict):
    char: str
    count: int

def get_num_words(text):
    words = text.split()
    return len(words)

def get_char_count(text: str) -> dict[str, int]:
    char_count = {}
    lower_case_text = text.lower()
    for char in lower_case_text:
        char_count[char] = char_count.get(char, 0) + 1
    return char_count

def sort_on(CharacterCount):
    return CharacterCount["count"]

def get_char_count_sorted(text: str) -> list[CharacterCount]:
    char_count = get_char_count(text)
    list_of_char_count = []

    for char, count in char_count.items():
        list_of_char_count.append(CharacterCount(char=char, count=count))

    list_of_char_count.sort(key=sort_on, reverse=True)

    for char_count_dict in list_of_char_count:
        if not char_count_dict["char"].isalpha():
            list_of_char_count.remove(char_count_dict)
            
    return list_of_char_count