class Solution:
    def lengthOfLongestSubstring(self, s: str) -> int:
        used=[]
        str=list(s)
        max=0
        for c in str:
            if c in used:
                if max < len(used):
                    max = len(used)
                del(used[:used.index(c)+1])
            used.append(c)

        if max < len(used):
            max = len(used)
        return max