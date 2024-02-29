class Solution:
    def wordPattern(self, pattern: str, s: str) -> bool:
        words=s.split() #O(n)
        if len(pattern)!=len(words): #O(1)
            return False

        char_to_word = {}
        word_to_char = {}

        for c, word in zip(pattern, words): #O?
            if c in char_to_word:
                if char_to_word[c] != word:
                    return False
            else:
                if word in word_to_char:
                    return False
                char_to_word[c] = word
                word_to_char[word] = c

        return True

