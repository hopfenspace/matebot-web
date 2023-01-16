export function formatAmount(amount: number): string {
    return new Intl.NumberFormat("de-DE", {
        currency: "EUR",
        style: "currency",
    }).format(amount / 100);
}
