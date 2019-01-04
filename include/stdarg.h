/* This header defines macros to enable variadic functions to be transpiled
 * with Leaven. It should be included instead of (or after) stdarg.h.
 */

#undef va_list
#undef va_start
#undef va_arg
#undef va_end

typedef void *leaven_va_list;

#define va_list leaven_va_list

void leaven_va_start(leaven_va_list *vl);
void *leaven_va_arg(leaven_va_list vl);

#define va_start(list, param) leaven_va_start(&list)
#define va_arg(list, type) (*(type *)leaven_va_arg(list))
#define va_end

