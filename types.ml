open Term
      
exception VarNotFound  
exception EquaNotFound 

type stype = 
  | VarS of string
  | AppS of stype * stype 
  | Nat
  | ListS of stype
  | ForAll of string * stype

type env = (string * stype) list            

type tequa = (stype * stype) list

type status =
  | Continue
  | Echec of string
  | Recommence
  | Fini

type typage_res = lterm * stype * status 

let rec stype_egal (st1 : stype) (st2 : stype) : (bool) =
  match (st1, st2) with
  | (VarS v1, VarS v2) -> v1 = v2
  | (AppS(st1, st2), AppS(st3, st4)) -> stype_egal st1 st3 && stype_egal st2 st4
  | (Nat, Nat) -> true 
  | (ListS st1, ListS st2) -> stype_egal st1 st2
  | (ForAll(x1, st1), ForAll(x2, st2)) -> x1 = x2 && stype_egal st1 st2 
  | _ -> false 
  
let equa_egal (eq1 : stype * stype) (eq2 : stype * stype) : (bool) =
  match (eq1, eq2) with
  | ((st1, st2), (st3, st4)) -> stype_egal st1 st3 && stype_egal st2 st4

let rec cherche_type (v : string) (e : env) : (stype) =
  match e with
  | [] -> raise VarNotFound
  | (v1, st)::q -> if v = v1 then st else cherche_type v q

let rec generalise (e : env) (st : stype) : (stype) =
  creer_for_all (vars_libres st) st
and creer_for_all (l : string list) (st : stype) : (stype)  =
  match l with
  | [] -> st
  | (h::q) -> ForAll(h, creer_for_all q st)
and vars_libres (st : stype) : (string list) =
  match st with 
  | VarS v -> [v]
  | AppS(st1, st2) -> fusionne (vars_libres st1) (vars_libres st2)
  | Nat -> []
  | ListS l -> vars_libres l
  | ForAll(v, l) -> enleve_var v (vars_libres l) 
and fusionne (l1 : string list) (l2 : string list) : (string list) =
  match l2 with
  | [] -> l1
  | (h::q) when contains h l1 -> fusionne l1 q
  | (h::q) -> fusionne (h::l1) q
and contains (str : string) (l : string list) : (bool) =
  match l with
  | [] -> false
  | (h::q) when h = str -> true
  | (h::q) -> contains str q
and enleve_var (str : string) (l : string list) : (string list) =
  match l with 
  | [] -> []
  | (h::q) when h = str -> enleve_var str q
  | (h::q) -> h::enleve_var str q 

let rec occur_check (str : string) (st : stype) : (bool) =
  match st with
  | VarS v -> str = v
  | AppS(arg, res) -> occur_check str arg || occur_check str res
  | ListS l -> occur_check str l
  | _ -> false

let rec substitue (str : string) (st : stype) (t : stype) : (stype) =
  match t with
  | VarS v -> if v = str then st else t
  | AppS(arg, res) -> AppS(substitue str st arg, substitue str st res)
  | Nat -> Nat
  | ListS l -> ListS (substitue str st l)
  | ForAll(x, l) -> ForAll(x, l)

let rec substitue_partout (str : string) (st : stype) (eqs : tequa) : (tequa) =
  match eqs with
  | [] -> []
  | (g, d)::q -> (substitue str st g, substitue str st d)::substitue_partout str st q 

let rec remove_tequa (eq : stype * stype) (eqs : tequa) : (tequa) =
  match eqs with
  | [] -> raise EquaNotFound
  | (h::q) -> if equa_egal eq h then q else h::remove_tequa eq q

let rec print_stype (st : stype) : (string) = 
  match st with
  | VarS v -> v
  | AppS(arg, res) -> "(" ^ print_stype arg ^ " -> " ^ print_stype res ^ ")"
  | Nat -> "Nat"
  | ListS l -> "[" ^ print_stype l ^ "]"
  | ForAll(x, l) -> (*"∀"*) "fa." ^ x ^ print_stype l 

let rec print_tequa (eq : tequa) : (string) =
  match eq with
  | [] -> ""
  | ((st1, st2)::q) -> print_stype st1 ^ " = " ^ print_stype st2 ^ "\n" ^ print_tequa q
