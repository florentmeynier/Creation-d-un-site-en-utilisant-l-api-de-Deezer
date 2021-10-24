open Term
open Types

type unif_res = tequa * status

let rec barendregt_stype (st : stype) (remp : (string * string) list) : (stype) =
  match st with
  | VarS v -> let res = contains v remp in
    if res = "" then st else VarS res 
  | AppS(st1, st2) -> AppS(barendregt_stype st1 remp, barendregt_stype st2 remp)
  | Nat -> Nat
  | ListS st1 -> ListS(barendregt_stype st1 remp)
  | ForAll(x, l) -> let nv = fresh_var () in
    ForAll(nv, barendregt_stype l ((x, nv)::remp))
and contains (str : string) (remp : (string * string) list) : (string) =
  match remp with
  | [] -> ""
  | ((v, s)::q) when v = str -> s
  | (h::q) -> contains str q 
and replace (str : string) (el : string * string) (remp : (string * string) list) : ((string * string) list) = 
  match remp with
  | [] -> [el]
  | ((x, _)::q) when x = str -> el::q
  | (h::q) -> h::replace str el q

let rec unification_etape (eqs : stype * stype) (e_final : tequa) : (unif_res) =
  match eqs with
    (*| (Nat, Nat) -> (remove_tequa (Nat, Nat) e_final, Continue)
  *)
  | (a1, a2) when stype_egal a1 a2 -> (remove_tequa  (a1, a2) e_final, Recommence) 
  | (ForAll(x, l), st) -> ((barendregt_stype l [], st)::(remove_tequa (ForAll(x, l), st) e_final), Recommence)
  | (st, ForAll(x, l)) -> ((barendregt_stype l [], st)::(remove_tequa (st, ForAll(x, l)) e_final), Recommence)
  | (VarS v1, st) when v1 = "but" -> (e_final, Continue)
  | (st, VarS v1) when v1 = "but" -> (e_final, Continue)
  | (VarS v1, VarS v2) -> (substitue_partout v1 (VarS v2) e_final, Recommence)
  | (VarS v, st) -> if occur_check v st
      then ([], Echec ("Occurence de " ^ v ^ " dans " ^ print_stype st)) 
      else (substitue_partout v st (remove_tequa (VarS v, st) e_final), Recommence)
  | (st, VarS v) -> if occur_check v st
      then ([], Echec ("Occurence de " ^ v ^ " dans " ^ print_stype st)) 
      else (substitue_partout v st (remove_tequa (st, VarS v) e_final), Recommence) 
  | (AppS(st1, st2), AppS(st3, st4)) -> ((st1, st3)::(st2, st4)::remove_tequa (AppS(st1, st2), AppS(st3, st4)) e_final, Continue)
  | (ListS l1, ListS l2) -> ((l1, l2)::(remove_tequa (ListS l1, ListS l2) e_final), Recommence)
  | _ -> ([], Echec "Flemme")
                                        
let rec unification (eq : tequa): (unif_res) = 
  unification_aux eq eq 
and unification_aux (eqs : tequa) (eqs2 : tequa) : (unif_res) =
  match eqs with
  | [] -> (eqs2, Fini)
  | (h::q) -> 
      (match unification_etape h eqs2 with
       | (res, Continue) -> unification_aux q res
       | (res, Recommence) -> unification_aux res res
       | (_, Echec s) -> ([], Echec s)
       | _ -> ([], Echec "Shouldn't happen")
      )

let rec print_unif_res (u : unif_res) : (string) =
  match u with
  | (res, _) -> aux res
and aux (st : (stype * stype) list) =
  match st with
  | [] -> ""
  | ((st1, st2)::q) -> print_stype st1 ^ print_stype st2 ^ aux q

