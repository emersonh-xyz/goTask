import clientPromise from "../../mongodb";

export async function GET() {
  try {
    const client = await clientPromise;
    const db = client.db("testdb");
    const result = await db.command({ ping: 1 });

    return new Response(
      JSON.stringify({ ok: result.ok === 1, message: "DB connection is good" }),
      { status: 200 }
    );
  } catch (error: any) {
    return new Response(
      JSON.stringify({ error: error.message }),
      { status: 500 }
    );
  }
}