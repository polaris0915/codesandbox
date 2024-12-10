
#include <iostream>
using namespace std;

const int N = 1000010;

int n, a[N];

void quick_sort(int* a, int left, int right){
   if (left >= right) return;

   int pivot = a[(left + right) / 2], i = left - 1, j = right + 1;
   while (i < j){
	   do i++; while (a[i] < pivot);
	   do j--; while (a[j] > pivot);
	   if (i < j) swap(a[i], a[j]);
   }

   quick_sort(a, left, j), quick_sort(a, j + 1, right);
}

int main(void){
   scanf("%d", &n);
   for (int i = 0; i < n; i++) scanf("%d", &a[i]);

   quick_sort(a, 0, n - 1);
   for (int i = 0; i < n; i++) printf("%d ", a[i]);
   printf("\n");
}
	