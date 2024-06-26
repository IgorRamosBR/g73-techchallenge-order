CREATE TABLE IF NOT EXISTS public.orders (
	"id" serial primary key,
	"coupon" text,
	"total_amount" numeric,
	"customer_cpf" text,
	"status" text,
	"created_at" timestamptz not null
);

CREATE TABLE IF NOT EXISTS public.order_items (
	"id" serial primary key,
	"order_id" int,
	"product_id" integer,
	"quantity" integer,
	"type" text,
	CONSTRAINT "FK_order_items_product" FOREIGN KEY (product_id) REFERENCES public.products(id),
	CONSTRAINT "FK_order" FOREIGN KEY (order_id) REFERENCES public.orders(id)
);