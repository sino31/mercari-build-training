
# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, x):
#         self.val = x
#         self.next = None

class Solution:
    def getIntersectionNode(self, headA: ListNode, headB: ListNode) -> Optional[ListNode]:
        alen,blen=0,0
        tmpA, tmpB=headA, headB
        while tmpA:
            alen+=1
            tmpA=tmpA.next
        while tmpB:
            blen+=1
            tmpB=tmpB.next
        while alen>blen:
            headA=headA.next
            alen-=1
        while blen>alen:
            headB=headB.next
            blen-=1

        while headA and headB:
            if headA==headB:
                return headA
            headA=headA.next
            headB=headB.next

        return None
