#### Variable syntax
```
access var name type;
private var score int;
score = 5;

score := 5 (assume private when using :=)
```

#### Function syntax
```
func modifiers <generics> Tr1, Tr2 f(a, b) {}
func public const <T implements Comparable> (T, int, error) num(T trait, string name) {
	return 22;
}
```
### Grammar

$$
\begin{align*}
	
	[\textcolor{red}{prog}] &\to [\textcolor{lime}{stmt}]^*
	\\
	[\textcolor{red}{stmt}] &\to \begin{cases}
		\textcolor{cyan}{var}\space\textcolor{yellow}{varIdent}\space[\textcolor{orange}{type}];\\
		\textcolor{yellow}{varIdent}=[\textcolor{lime}{expr}];\\
		*\textcolor{yellow}{varIdent}=[\textcolor{lime}{expr}];\\
		[\textcolor{lime}{scope}]\\
		[\textcolor{lime}{if}]\\
		\textcolor{cyan}{while}([\textcolor{orange}{bool}\textcolor{lime}{expr}])[\textcolor{lime}{scope}]\\
		\textcolor{cyan}{break};\\
		\textcolor{cyan}{continue};\\
		\textcolor{cyan}{func}\space\text{intLiteral}\space\textcolor{yellow}{funcIdent}(\textcolor{yellow}{param1}\space[\textcolor{orange}{type}],^*)[\textcolor{lime}{scope}]\\
		[\textcolor{lime}{funcCall}];\\		
		\textcolor{cyan}{return}\space[\textcolor{lime}{expr}],^*;\\
		\textcolor{cyan}{syscall}([\textcolor{lime}{expr}],^*);\\
	\end{cases}
	\\
	[\textcolor{red}{expr}] &\to \begin{cases}
		[\textcolor{orange}{int}\textcolor{lime}{expr}]\\
		[\textcolor{orange}{bool}\textcolor{lime}{expr}]\\
	\end{cases}
	\\
	[\textcolor{orange}{int}\textcolor{red}{Expr}] &\to \begin{cases}
		[\textcolor{orange}{int}\textcolor{lime}{term}]\\
		[\textcolor{orange}{int}\textcolor{lime}{binExpr}]\\
	\end{cases}
	\\
	[\textcolor{orange}{int}\textcolor{red}{BinExpr}] &\to \begin{cases}
		[\textcolor{orange}{int}\textcolor{lime}{expr}]\%[\textcolor{orange}{int}\textcolor{lime}{expr}] & \textcolor{magenta}{prec=1}\\
		[\textcolor{orange}{int}\textcolor{lime}{expr}]*[\textcolor{orange}{int}\textcolor{lime}{expr}] & \textcolor{magenta}{prec=1}\\
		[\textcolor{orange}{int}\textcolor{lime}{expr}]/[\textcolor{orange}{int}\textcolor{lime}{expr}] & \textcolor{magenta}{prec=1}\\
		[\textcolor{orange}{int}\textcolor{lime}{expr}]+[\textcolor{orange}{int}\textcolor{lime}{expr}] & \textcolor{magenta}{prec=0}\\
		[\textcolor{orange}{int}\textcolor{lime}{expr}]-[\textcolor{orange}{int}\textcolor{lime}{expr}] & \textcolor{magenta}{prec=0}\\
	\end{cases}
	\\
	[\textcolor{orange}{int}\textcolor{red}{Term}] &\to \begin{cases}
		-[\textcolor{orange}{int}\textcolor{lime}{term}]\\
		\text{intLiteral}\\
		\textcolor{yellow}{varIdent}\\
		([\textcolor{orange}{int}\textcolor{lime}{expr}])\\
		[\textcolor{lime}{funcCall}]\\
		\&\textcolor{yellow}{varIdent}&&&\text{(need to move this to separate pointer type)}\\
		*\textcolor{yellow}{varIdent}\\ 
	\end{cases}
	\\
	[\textcolor{orange}{bool}\textcolor{red}{Expr}] &\to \begin{cases}
		[\textcolor{orange}{bool}\textcolor{lime}{term}]\\
		[\textcolor{orange}{bool}\textcolor{lime}{binExpr}]\\
	\end{cases}
	\\
	[\textcolor{orange}{bool}\textcolor{red}{BinExpr}] &\to \begin{cases}
		[\textcolor{orange}{bool}\textcolor{lime}{expr}]\&\&[\textcolor{orange}{cool}\textcolor{lime}{expr}] & \textcolor{magenta}{prec=1}\\
		[\textcolor{orange}{bool}\textcolor{lime}{expr}]||[\textcolor{orange}{bool}\textcolor{lime}{expr}] & \textcolor{magenta}{prec=0}\\
	\end{cases}
	\\
	[\textcolor{orange}{bool}\textcolor{red}{Term}] &\to \begin{cases}
		![\textcolor{orange}{bool}\textcolor{lime}{term}]\\
		\textcolor{yellow}{varIdent}\\
		([\textcolor{orange}{bool}\textcolor{lime}{expr}])\\
		[\textcolor{lime}{funcCall}]\\
		*\textcolor{yellow}{varIdent}\\
		[\textcolor{orange}{bool}\textcolor{lime}{term}] [\textcolor{lime}{relOp}] [\textcolor{orange}{bool}\textcolor{lime}{term}]\\
		[\textcolor{orange}{int}\textcolor{lime}{expr}] [\textcolor{lime}{relOp}] [\textcolor{orange}{int}\textcolor{lime}{expr}]\\
	\end{cases}
	\\
	[\textcolor{red}{relOp}] &\to \begin{cases}
		==\\
		!=\\
		<\\
		>\\
		>=\\
		<=\\
	\end{cases}
	\\
	[\textcolor{red}{scope}] &\to \{[\textcolor{lime}{stmt}]^*\}
	\\
	[\textcolor{red}{if}] &\to \textcolor{cyan}{if}([\textcolor{orange}{bool}\textcolor{lime}{expr}])[\textcolor{lime}{scope}]<\textcolor{cyan}{else}\space[\textcolor{lime}{else}]>\\
	[\textcolor{red}{else}] &\to \begin{cases}
		[\textcolor{lime}{if}]\\
		[\textcolor{lime}{scope}]\\
	\end{cases}\\

	[\textcolor{red}{funcCall}] &\to \textcolor{yellow}{funcIdent}([\textcolor{lime}{expr}],^*)\\

\end{align*}
$$

$$
\begin{align*}

	[\textcolor{orange}{type}] &\to \begin{cases}
		[\textcolor{orange}{baseType}]\\
		*[\textcolor{orange}{type}]\\
	\end{cases}
	\\
	[\textcolor{orange}{baseType}] &\to \begin{cases}
		\textcolor{cyan}{bool}\\
		\textcolor{cyan}{int}\\
		\textcolor{cyan}{char}\\
	\end{cases}
	\\

\end{align*}
$$

