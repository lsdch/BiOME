import { Meta, User } from "@/api";



export function isOwner<
  Item extends { meta?: Meta }
>(user: User, item: Item) {
  return item.meta?.created_by?.id === user.id
}