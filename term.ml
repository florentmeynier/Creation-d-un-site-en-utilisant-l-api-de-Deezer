type lterm = 
  | Var of string 
  | App of lterm * lterm
  | Abs of string * lterm
  | N of int
  | Add of lterm * lterm
  | Sub of lterm * lterm
  | List of lterm liste
  | Head of lterm
	| Tail of lterm
  | IfZero of lterm * lterm * lterm
  | IfEmpty of lterm * lterm * lterm
  | Fix 
  | Let of string * lterm * lterm
and 'a liste = Cons of 'a * 'a liste | Nil

let cvar (v : string) : (lterm) = Var v

let capp (fpos : lterm) (apos : lterm) : (lterm) = App(fpos, apos)

let cabs (v : string) (corps : lterm) : (lterm) = Abs(v, corps)

let cint (n : int) : (lterm) = N n

let cadd (g : lterm) (d : lterm) : (lterm) = Add(g, d)

let csub (g : lterm) (d : lterm) : (lterm) = Sub(g, d)

let rec clist (l : lterm list) : (lterm) =
	List (clist_rec l)
and clist_rec (l : lterm list) : (lterm liste) =
	match l with
	| [] -> Nil
	| (h::q) -> Cons(h, clist_rec q)

let chead (lt : lterm) : (lterm) = Head lt

let ctail (lt : lterm) : (lterm) = Tail lt

let cizte (cond : lterm) (t : lterm) (e : lterm) : (lterm) = IfZero(cond, t, e)

let iete (cond : lterm) (t : lterm) (e : lterm) : (lterm) = IfEmpty(cond, t, e)

let cfix () : (lterm) = Fix

let clet (v : string) (lt1 : lterm) (lt2 : lterm) : (lterm) = Let(v, lt1, lt2)

let counter = ref 0

let fresh_var() : (string) =
  counter := !counter + 1;
  "T" ^ string_of_int(!counter) 

let rec barendregt_rec_lterm (lt : lterm) (remp : (string * string) list) : (lterm) =
  match lt with 
  | Var v -> let res = contains v remp in
    if res = "" then lt else Var res
  | App(lt1, lt2) -> App(barendregt_rec_lterm lt1 remp, barendregt_rec_lterm lt2 remp)
  | Abs(x, lt1) -> let nv = fresh_var () in
    Abs(nv, barendregt_rec_lterm lt1 (replace x (x, nv) remp))
  | List l -> List (barendregt_rec_list l remp)
  | IfZero(corps, fpos, apos) -> IfZero(barendregt_rec_lterm corps remp, barendregt_rec_lterm fpos remp, barendregt_rec_lterm apos remp)
  | IfEmpty(corps, fpos, apos) -> IfEmpty(barendregt_rec_lterm corps remp, barendregt_rec_lterm fpos remp, barendregt_rec_lterm apos remp)
  | Let(v, fpos, apos) -> let nv = fresh_var () in
    Let(nv, barendregt_rec_lterm fpos remp, barendregt_rec_lterm apos (replace v (v, nv) remp))
  | _ -> lt
and contains (str : string) (remp : (string * string) list) : string =
  match remp with
  | [] -> ""
  | ((v, s)::q) when v = str -> s
  | (h::q) -> contains str q 
and replace (str : string) (el : string * string) (remp : (string * string) list) : ((string * string) list) = 
  match remp with
  | [] -> [el]
  | ((x, _)::q) when x = str -> el::q
  | (h::q) -> h::replace str el q
and barendregt_rec_list (l : lterm liste) (remp : (string * string) list) : (lterm liste) =
  match l with
  | Nil -> Nil
  | Cons(e, q) -> Cons(barendregt_rec_lterm e remp, barendregt_rec_list q remp)

let barendregt_lterm (lt : lterm) : (lterm) =
  barendregt_rec_lterm lt []

let rec instancie (lt : lterm) (x : string) (a : lterm) : (lterm) =
  match lt with
  | Var v when v = x -> a
  | App(lt1, lt2) -> App(instancie lt1 x a, instancie lt2 x a)
  | Abs(v, lt1) -> Abs(v, instancie lt1 x a)
  | List l -> List(instancie_list l x a)
  | IfZero(corps, fpos, apos) -> IfZero(instancie corps x a, instancie fpos x a, instancie apos x a)
  | IfEmpty(corps, fpos, apos) -> IfEmpty(instancie corps x a, instancie fpos x a, instancie apos x a)
  | _ -> lt
and instancie_list (l : lterm liste) (x : string) (a : lterm) : (lterm liste) =
  match l with
  | Nil -> Nil
  | Cons(e, q) -> Cons(instancie e x a, instancie_list q x a)

let rec print_lterm (lt : lterm) : (string) =
  match lt with
  | Var v -> v
  | App(fpos, apos) -> "(" ^ print_lterm fpos ^ " " ^ print_lterm apos ^ ")"
  | Abs(v, t) ->  "(fun " ^ v ^ " -> " ^ print_lterm t ^ ")"
  | N n -> string_of_int n
  | Add(g, d) -> "(" ^ print_lterm g ^ " + " ^ print_lterm d ^ ")"
  | Sub(g, d) -> "(" ^ print_lterm g ^ " - " ^ print_lterm d ^ ")"
  | List l -> "[" ^ print_list l
  | Head l ->  "(hd" ^ print_lterm l ^ ")"
  | Tail l -> "tl" ^ print_lterm l ^ ")"
  | IfZero(corps, fpos, apos) -> "if0 " ^ print_lterm corps ^ " then " ^ print_lterm fpos ^ " else " ^ print_lterm apos
  | IfEmpty(corps, fpos, apos) -> "ifempty " ^ print_lterm corps ^ " then " ^ print_lterm fpos ^ " else " ^ print_lterm apos
  | Fix -> "fix"
  | Let(v, fpos, apos) -> "(let " ^ v ^ " = " ^ print_lterm fpos ^ " in " ^ print_lterm apos ^ ")"
and print_list (l : lterm liste) : (string) =
  match l with
  | Nil -> "]"
  | Cons(e, Nil) -> print_lterm e ^ "]"
  | Cons(e, q) -> print_lterm e ^ "," ^ print_list q