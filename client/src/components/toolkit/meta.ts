import { Meta, User } from "@/api";



export function isOwner<
  Item extends { id: string, meta?: Meta }
>(user: User, item: Item) {
  return item.meta?.created_by_user.id === user.id
}